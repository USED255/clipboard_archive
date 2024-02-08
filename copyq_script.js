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

function main() {
    if (hasBigData()) {
        return;
    }
    url = url + str(Item["application/x-copyq-user-copy-time"]);
    Item = JSON.stringify({"Data": toBase64(pack(Item))});
    networkPost(url, Item).data;
}

main()
