const https = require('https');
const fs = require('fs');
const path = require('path');

const fonts = [
  {
    name: 'Poppins-Bold.ttf',
    url: 'https://raw.githubusercontent.com/google/fonts/main/ofl/poppins/Poppins-Bold.ttf'
  },
  {
    name: 'Poppins-SemiBold.ttf',
    url: 'https://raw.githubusercontent.com/google/fonts/main/ofl/poppins/Poppins-SemiBold.ttf'
  },
  {
    name: 'Poppins-Medium.ttf',
    url: 'https://raw.githubusercontent.com/google/fonts/main/ofl/poppins/Poppins-Medium.ttf'
  },
  {
    name: 'Inter-Regular.ttf',
    url: 'https://raw.githubusercontent.com/google/fonts/main/ofl/inter/Inter[slnt,wght].ttf'
  },
  {
    name: 'Inter-Medium.ttf',
    url: 'https://raw.githubusercontent.com/google/fonts/main/ofl/inter/Inter[slnt,wght].ttf'
  }
];

const fontsDir = path.join(__dirname, '../assets/fonts');

// Create fonts directory if it doesn't exist
if (!fs.existsSync(fontsDir)) {
  fs.mkdirSync(fontsDir, { recursive: true });
}

fonts.forEach(font => {
  const filePath = path.join(fontsDir, font.name);
  const file = fs.createWriteStream(filePath);

  https.get(font.url, response => {
    response.pipe(file);
    file.on('finish', () => {
      file.close();
      console.log(`Downloaded ${font.name}`);
    });
  }).on('error', err => {
    fs.unlink(filePath);
    console.error(`Error downloading ${font.name}:`, err.message);
  });
}); 