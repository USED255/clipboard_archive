package database

const CreateItemsTableQuery = `
CREATE TABLE "clipboard_items" (
	"id" integer,
	"created_at" datetime,
	"updated_at" datetime,
	"deleted_at" datetime,
	"clipboard_item_time" integer UNIQUE,
	"clipboard_item_text" text,
	"clipboard_item_hash" text UNIQUE,
	"clipboard_item_data" text,
	PRIMARY KEY ("id")
);
INSERT INTO "main"."clipboard_items" (
	"id", 
	"created_at", 
	"updated_at", 
	"deleted_at", 
	"clipboard_item_time", 
	"clipboard_item_text", 
	"clipboard_item_hash", 
	"clipboard_item_data"
) 
VALUES (
	"499", 
	"2022-03-13 13:22:43.238644233+08:00", 
	"2022-03-13 13:22:43.238644233+08:00", 
	"", 
	"1647146952858", 
	"migrate", 
	"2cb5fed12b27c377de172eb922161838b1343adf55Dbd9Db39aa50391f1fc2c7", 
	"/////gAAAAAAA="
);
`

const CreateConfigsTableQuery = `
CREATE TABLE "configs" (
	"key" text,
	"value" text,
	PRIMARY KEY ("key")
);
`
