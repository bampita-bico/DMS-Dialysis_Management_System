/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Darker gray palette for African context
        gray: {
          50: '#f8f9fa',   // Slightly darker than default
          100: '#e9ecef',  // Reduced brightness
          200: '#dee2e6',
          300: '#ced4da',
          400: '#adb5bd',
          500: '#6c757d',
          600: '#495057',
          700: '#343a40',
          800: '#212529',
          900: '#181c1f',
        },
      },
    },
  },
  plugins: [],
}
