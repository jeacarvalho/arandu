module.exports = {
  theme: {
    extend: {
      colors: {
        arandu: {
          primary: '#1E3A5F',
          secondary: '#3A7D6B',
          insight: '#D4A84F',
          background: '#F7F8FA',
          text: '#1F2937'
        }
      },
      fontFamily: {
        // Interface: Inter (limpa e funcional)
        'sans': ['Inter', 'ui-sans-serif', 'system-ui'],
        // Conteúdo Clínico: Source Serif (conforto para leitura e escrita)
        'serif': ['Source Serif 4', 'Source Serif Pro', 'ui-serif', 'Georgia'],
      }
    }
  }
}