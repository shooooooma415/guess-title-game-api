-- Remove sample themes

DELETE FROM themes WHERE title IN (
  '夏祭り',
  'コーヒー',
  '富士山',
  'ラーメン',
  'サッカー',
  'お寿司',
  '桜',
  '温泉',
  '花火',
  'アニメ'
);
