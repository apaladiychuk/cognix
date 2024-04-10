
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
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        success: {
          lighter: 'hsl(var(--success-lighter))',
          light: 'hsl(var(--success-light))',
          DEFAULT: 'hsl(var(--success))',
          dark: 'hsl(var(--success-dark))',
          darker: 'hsl(var(--success-darker))',
        },
        destructive: {
          foreground: 'hsl(var(--destructive-foreground))',
          lighter: 'hsl(var(--destructive-lighter))',
          light: 'hsl(var(--destructive-light))',
          DEFAULT: 'hsl(var(--destructive))',
          dark: 'hsl(var(--destructive-dark))',
          darker: 'hsl(var(--destructive-darker))',
        },
        'conversation-status-open': {
          DEFAULT: 'hsl(var(--conversation-status-open))',
        },
        'conversation-status-pending': {
          DEFAULT: 'hsl(var(--conversation-status-pending))',
        },
        'conversation-status-overdue': {
          DEFAULT: 'hsl(var(--conversation-status-overdue))',
        },
        'conversation-status-resolved': {
          DEFAULT: 'hsl(var(--conversation-status-resolved))',
        },
        'conversation-status-sneakpeek': {
          DEFAULT: 'hsl(var(--conversation-status-sneakpeek))',
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
        '2sm': '0.75rem',
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
