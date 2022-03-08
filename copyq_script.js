copyq:

var minBytes = 250*1000
function hasBigData() {
    var itemSize = 0
    var formats = dataFormats()
    for (var i in formats) {
        itemSize += data(formats[i]).size()
        if (itemSize >= minBytes)
            return true
    }
    return false
}

function main(){
    if (hasBigData()) {
        return;
    }
    ClipboardItem = getItem(0);
    ClipboardItemData = toBase64(pack(ClipboardItem));
    ClipboardItemHash = sha256sum(ClipboardItemData);
    ClipboardItemTime = parseInt(str(ClipboardItem["application/x-copyq-user-copy-time"]));
    ClipboardItemText = str(ClipboardItem[mimeText]);
    ClipboardItemObject = {"ClipboardItemTime":ClipboardItemTime, "ClipboardItemText":ClipboardItemText, "ClipboardItemHash":ClipboardItemHash, "ClipboardItemData":ClipboardItemData};
    ClipboardItemJson = JSON.stringify(ClipboardItemObject);
    url = "http://localhost:8080/api/v1/ClipboardItem";
    networkPost(url,ClipboardItemJson).data;
}
main()