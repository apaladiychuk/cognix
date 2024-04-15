
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ['class'],
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      colors: {
        border: 'var(--border)',
        input: 'var(--input)',
        ring: 'var(--ring)',
        background: 'var(--background)',
        foreground: 'var(--foreground)',
        primary: {
          DEFAULT: 'var(--primary)',
          foreground: 'var(--primary-foreground)',
        },
        secondary: {
          DEFAULT: 'var(--secondary)',
          foreground: 'var(--secondary-foreground)',
        },
        muted: {
          DEFAULT: 'var(--muted)',
          foreground: 'var(--muted-foreground)',
        },
        accent: {
          DEFAULT: 'var(--accent)',
          foreground: 'var(--accent-foreground)',
        },
        popover: {
          DEFAULT: 'var(--popover)',
          foreground: 'var(--popover-foreground)',
        },
        card: {
          DEFAULT: 'var(--card)',
          foreground: 'var(--card-foreground)',
        },
        success: {
          lighter: 'var(--success-lighter)',
          light: 'var(--success-light)',
          DEFAULT: 'var(--success)',
          dark: 'var(--success-dark)',
          darker: 'var(--success-darker)',
        },
        destructive: {
          foreground: 'var(--destructive-foreground)',
          lighter: 'var(--destructive-lighter)',
          light: 'var(--destructive-light)',
          DEFAULT: 'var(--destructive)',
          dark: 'var(--destructive-dark)',
          darker: 'var(--destructive-darker)',
        },
        'conversation-status-open': {
          DEFAULT: 'var(--conversation-status-open)',
        },
        'conversation-status-pending': {
          DEFAULT: 'var(--conversation-status-pending)',
        },
        'conversation-status-overdue': {
          DEFAULT: 'var(--conversation-status-overdue)',
        },
        'conversation-status-resolved': {
          DEFAULT: 'var(--conversation-status-resolved)',
        },
        'conversation-status-sneakpeek': {
          DEFAULT: 'var(--conversation-status-sneakpeek)',
        },
      },
      borderRadius: {
        lg: `var(--radius)`,
        md: `calc(var(--radius) - 2px)`,
        sm: 'calc(var(--radius) - 4px)',
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
      },
      fontSize: {
        '2sm': '0.85rem',
        '2xs': '0.625rem',
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
      },
    },
  },
  plugins: [],
};
