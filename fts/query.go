package fts

const createFts5TableQuery = `
	CREATE VIRTUAL TABLE clipboard_items_fts USING fts5(
		clipboard_item_time, 
		clipboard_item_text, 
		content = clipboard_items, 
		content_rowid = clipboard_item_time
	);
	
	CREATE TRIGGER clipboard_items_ai AFTER INSERT ON clipboard_items BEGIN
		INSERT INTO clipboard_items_fts(
			rowid, 
			clipboard_item_text
		) 
		VALUES (
			new.clipboard_item_time, 
			new.clipboard_item_text
		);
	END;
		
	CREATE TRIGGER clipboard_items_ad AFTER DELETE ON clipboard_items BEGIN
		INSERT INTO clipboard_items_fts(
			clipboard_items_fts, 
			rowid, 
			clipboard_item_text
		) 
		VALUES(
			"delete", 
			old.clipboard_item_time, 
			old.clipboard_item_text
		);
	END;
		
	CREATE TRIGGER clipboard_items_au AFTER UPDATE ON clipboard_items BEGIN
		INSERT INTO clipboard_items_fts(
			clipboard_items_fts, 
			rowid, 
			clipboard_item_text
		) 
		VALUES(
			"delete", 
			old.clipboard_item_time, 
			old.clipboard_item_text
		);
		INSERT INTO clipboard_items_fts(
			rowid, 
			clipboard_item_text
		) 
		VALUES (
			new.clipboard_item_time, 
			new.clipboard_item_text
		);
	END;
`

const insertFts5TableQuery = `
INSERT INTO clipboard_items_fts (
	rowid, 
	clipboard_item_text
)
SELECT clipboard_items.clipboard_item_time, clipboard_items.clipboard_item_text 
FROM clipboard_items;
`
