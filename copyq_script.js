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
    x_item = getItem(0);
    x_data = toBase64(pack(x_item));
    x_hash = sha256sum(x_data);
    x_time = parseInt(str(x_item["application/x-copyq-user-copy-time"]));
    x_text = str(x_item[mimeText]);
    x_object = {"x_time": x_time, "x_text": x_text, "x_hash": x_hash,"x_data": x_data};
    x_json = JSON.stringify(x_object);
    url = "http://localhost:8088/api/v1/items";
    networkPost(url,x_json).data;
}
main()