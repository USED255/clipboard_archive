// https://copyq.readthedocs.io/en/latest/scripting-api.html
copyq:

var minBytes = 250 * 1000;
var url = "https://127.0.0.1:8080/api/v2/Item";

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

function NewItem(Item) {
    return JSON.stringify({
        "Time": parseInt(str(Item["application/x-copyq-user-copy-time"])),
        "Data": toBase64(pack(Item)),
    });
}

function main() {
    if (hasBigData()) {
        return;
    }
    Item = NewItem(getItem(0));
    networkPost(url, Item).data;
}

main()
