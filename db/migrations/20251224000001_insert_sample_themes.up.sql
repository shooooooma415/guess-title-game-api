-- Insert sample themes for development and testing

INSERT INTO themes (id, title, hint) VALUES
  (gen_random_uuid(), '日本映画', '日本で制作された映画作品'),
  (gen_random_uuid(), '日本アニメ', '日本のアニメーション作品'),
  (gen_random_uuid(), '日本の有名人', '日本で有名な人物'),
  (gen_random_uuid(), '日本料理', '日本の伝統的な食べ物'),
  (gen_random_uuid(), '地名', '日本や世界の場所の名前');
