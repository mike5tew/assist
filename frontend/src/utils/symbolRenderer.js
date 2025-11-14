export const renderMedicalText = (text, useOriginal = true) => {
  if (!text) return '';

  // If original text with symbols is available, use it
  // Otherwise, convert normalized text back to symbols for display
  const symbolMap = {
    'mu': 'μ',
    'alpha': 'α',
    'beta': 'β',
    'gamma': 'γ',
    'delta': 'δ',
    'microliter': 'µl',
    'microgram': 'µg',
    'degrees-Celsius': '°C',
  };

  let displayText = text;
  
  for (const [normalized, symbol] of Object.entries(symbolMap)) {
    const regex = new RegExp(`\\b${normalized}\\b`, 'gi');
    displayText = displayText.replace(regex, symbol);
  }

  return displayText;
};
