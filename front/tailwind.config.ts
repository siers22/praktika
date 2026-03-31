import type { Config } from 'tailwindcss'

const config: Config = {
  content: ['./src/**/*.{js,ts,jsx,tsx,mdx}'],
  theme: {
    extend: {
      colors: {
        primary: '#006666',
        'primary-dark': '#004d4d',
        'primary-light': '#008080',
        secondary: '#F1F2F5',
        surface: '#E7E5E4',
        'surface-dark': '#d4d1cf',
        'text-main': '#1E2938',
        'text-muted': '#6B7280',
        success: '#00A63D',
        warning: '#FE9900',
        danger: '#FF2157',
        info: '#0066CC',
      },
      fontFamily: {
        mono: ['"Space Mono"', '"JetBrains Mono"', 'monospace'],
      },
      boxShadow: {
        neu: '6px 6px 12px #c8c6c4, -6px -6px 12px #ffffff',
        'neu-sm': '3px 3px 6px #c8c6c4, -3px -3px 6px #ffffff',
        'neu-lg': '10px 10px 20px #c8c6c4, -10px -10px 20px #ffffff',
        'neu-inset': 'inset 4px 4px 8px #c8c6c4, inset -4px -4px 8px #ffffff',
        'neu-inset-sm': 'inset 2px 2px 5px #c8c6c4, inset -2px -2px 5px #ffffff',
        'neu-primary': '6px 6px 12px #004d4d33, -6px -6px 12px #ffffff',
      },
      borderRadius: {
        neu: '16px',
        'neu-sm': '10px',
        'neu-lg': '24px',
      },
    },
  },
  plugins: [],
}
export default config
