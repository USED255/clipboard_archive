// https://copyq.readthedocs.io/en/latest/scripting-api.html
copyq:

var minBytes = 250 * 1000;
var url = "https://127.0.0.1:8080/api/v1/ClipboardItem";

function hasBigData() {
    var itemSize = 0;
    var formats = dataFormats();
    for (var i in formats) {
        itemSize += data(formats[i]).size();
        if (itemSize >= minBytes)
            return true
    }
    return false
}

function clipboardItem(Item) {
    ClipboardItemData = toBase64(pack(Item));
    ClipboardItemHash = sha256sum(ClipboardItemData);
    ClipboardItemTime = parseInt(str(Item["application/x-copyq-user-copy-time"]));
    ClipboardItemText = str(Item[mimeText]);
    ClipboardItemObject = {
        "ClipboardItemTime": ClipboardItemTime,
        "ClipboardItemText": ClipboardItemText,
        "ClipboardItemHash": ClipboardItemHash,
        "ClipboardItemData": ClipboardItemData
    };
    return JSON.stringify(ClipboardItemObject);
}

function main() {
    if (hasBigData()) {
        return;
    }
    ClipboardItem = clipboardItem(getItem(0));
    networkPost(url, ClipboardItem).data;
}

main()
