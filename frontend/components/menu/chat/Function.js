export function getFileExtensionData(filename){
    let regex = /(?:\.([^.]+))?$/;
    let matches = regex.exec(filename);
    let extension = matches ? matches[1] : '';

    switch(extension) {
        case 'png':
            return 'data:image/png;base64,'
        case 'jpeg' || 'jpg':
            return 'data:image/jpeg;base64,'
    }

    return 'data:image/jpeg;base64,'
}

export function generateFileName(filename, length) {
    let extension = getFileExtension(filename)

    let result = '';
    let characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let charactersLength = characters.length;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }

    return result + '.' + extension;
}

function getFileExtension(filename){
    let regex = /(?:\.([^.]+))?$/;

    let matches = regex.exec(filename);
    let extension = matches ? matches[1] : '';

    return extension
}

export function GetStringMonth(month) {
    switch(month) {
        case 0:
            return 'Января'
        case 1:
            return 'Феврлаля'
        case 2:
            return 'Марта'
        case 3:
            return 'Апреля'
        case 4:
            return 'Мая'
        case 5:
            return 'Июня'
        case 6:
            return 'Июля'
        case 7:
            return 'Августа'
        case 8:
            return 'Сентября'
        case 9:
            return 'Октября'
        case 10:
            return 'Ноября'
        case 11:
            return 'Декабря'
    }
}