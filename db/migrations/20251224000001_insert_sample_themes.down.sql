-- Remove sample themes

DELETE FROM themes WHERE title IN (
  '日本映画',
  '日本アニメ',
  '日本の有名人',
  '日本料理',
  '地名'
);
