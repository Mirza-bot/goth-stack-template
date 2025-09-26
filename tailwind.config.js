/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "internal/templates/**/*.templ" ],
  theme: { extend: {} },
  plugins: [],
  daisyui: {
    themes: ["nord", "sunset"],
    base: true,
    styled: true,
    utils: true,
    logs: false
  }
}
