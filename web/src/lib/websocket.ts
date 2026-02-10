'use client';

import { useEffect, useRef, useCallback } from 'react';
import { WSEvent } from './types';
import { api } from './api';

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/api/ws';

export function useWebSocket(onMessage: (event: WSEvent) => void) {
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);
    const reconnectAttemptsRef = useRef(0);
    const onMessageRef = useRef(onMessage);

    // Keep the callback ref up to date
    useEffect(() => {
        onMessageRef.current = onMessage;
    }, [onMessage]);

    const connect = useCallback(() => {
        const token = api.getToken();
        if (!token) return;

        // If we already have a connection (OPEN or CONNECTING), don't reconnect
        if (wsRef.current && (wsRef.current.readyState === WebSocket.OPEN || wsRef.current.readyState === WebSocket.CONNECTING)) {
            return;
        }

        try {
            // Close existing connection ONLY if it's not already closed
            if (wsRef.current) {
                wsRef.current.onclose = null;
                wsRef.current.onerror = null;
                wsRef.current.onmessage = null;
                wsRef.current.onopen = null;
                wsRef.current.close();
                wsRef.current = null;
            }

            console.log('Connecting to WebSocket...');
            const ws = new WebSocket(`${WS_URL}?token=${token}`);

            ws.onopen = () => {
                console.log('WebSocket connected');
                reconnectAttemptsRef.current = 0;
            };

            ws.onmessage = (event) => {
                try {
                    // Handle batching (messages joined with \n)
                    const lines = event.data.split('\n');
                    for (const line of lines) {
                        if (!line.trim()) continue;
                        try {
                            const data = JSON.parse(line) as WSEvent;
                            if (onMessageRef.current) {
                                onMessageRef.current(data);
                            }
                        } catch (err) {
                            console.warn('Failed to parse WebSocket message line:', err);
                        }
                    }
                } catch (err) {
                    console.error('WebSocket message processing error:', err);
                }
            };

            ws.onerror = () => {
                // Don't log full error object as it's often empty/unhelpful in browser
                console.warn('WebSocket error occurred');
            };

            ws.onclose = (event) => {
                wsRef.current = null;

                // Don't reconnect if it was a normal close or if we are unmounted
                if (event.wasClean) {
                    console.log('WebSocket closed normally');
                    return;
                }

                console.log(`WebSocket disconnected: ${event.code} ${event.reason}`);

                // Exponential backoff reconnection
                const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
                reconnectAttemptsRef.current++;
                reconnectTimeoutRef.current = setTimeout(connect, delay);
            };

            wsRef.current = ws;
        } catch (err) {
            console.error('Failed to create WebSocket connection:', err);
        }
    }, []); // No dependencies - connect logic is static

    const disconnect = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
        }
        if (wsRef.current) {
            // Prevent auto-reconnect on manual disconnect
            wsRef.current.onclose = null;
            wsRef.current.close();
            wsRef.current = null;
        }
    }, []);

    const send = useCallback((event: string, data: unknown) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({ event, data }));
        }
    }, []);

    useEffect(() => {
        connect();
        return () => disconnect();
    }, [connect, disconnect]);

    return { send, reconnect: connect, disconnect };
}
