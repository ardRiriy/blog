-- 修正後の内容: 存在確認を追加
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'articles' AND column_name = 'tags') THEN
        ALTER TABLE articles DROP COLUMN tags;
    END IF;
END $$;

DROP INDEX IF EXISTS idx_articles_tags;

ALTER TABLE articles ADD COLUMN subtitle TEXT;
