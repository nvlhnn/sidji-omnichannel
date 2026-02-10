import type { Config } from 'tailwindcss';

const config: Config = {
    content: [
        './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
        './src/components/**/*.{js,ts,jsx,tsx,mdx}',
        './src/app/**/*.{js,ts,jsx,tsx,mdx}',
    ],
    theme: {
        extend: {
            backgroundImage: {
                'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
                'gradient-conic':
                    'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
            },
            colors: {
                primary: {
                    DEFAULT: 'var(--primary)',
                    hover: 'var(--primary-hover)',
                    light: 'var(--primary-light)',
                },
                background: {
                    DEFAULT: 'var(--background)',
                    secondary: 'var(--background-secondary)',
                    tertiary: 'var(--background-tertiary)',
                },
                foreground: {
                    DEFAULT: 'var(--foreground)',
                    secondary: 'var(--foreground-secondary)',
                    muted: 'var(--foreground-muted)',
                },
            },
        },
    },
    plugins: [],
};
export default config;
