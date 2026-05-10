-- Insert sample themes for development and testing

INSERT INTO themes (id, title, hint) VALUES
  (gen_random_uuid(), 'となりのトトロ', '日本の映画'),
  (gen_random_uuid(), 'ワンピース', '日本のアニメ'),
  (gen_random_uuid(), '大谷翔平', '日本の有名人'),
  (gen_random_uuid(), 'おでん', '日本料理'),
  (gen_random_uuid(), '大阪城', '地名');
