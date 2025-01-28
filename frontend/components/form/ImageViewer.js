import { Image } from 'react-native';

let sizeImg = 40
let sizeImgContact = 50
let sizeImgChat = 50

export default function ImageViewer({ placeholderImageSource, selectedImage, size, chat}) {
    const imageSource = selectedImage && selectedImage !== '' ? selectedImage : placeholderImageSource;

    if (imageSource === 2 || imageSource === 4) {
      return
    }

    if (chat !== undefined) {
      return <Image
              source={{uri: imageSource}}
              style={{
                width: sizeImgChat,
                height: sizeImgChat,
                borderRadius: sizeImgChat,
              }} 
            />;
    }

    if (size !== undefined) {
      return <Image
              source={{uri: imageSource}}
              style={{
                width: sizeImgContact,
                height: sizeImgContact,
                borderRadius: sizeImgContact,
              }} 
            />;
    }

    return <Image
            source={imageSource}
            style={{
              width: sizeImg,
              height: sizeImg,
              borderRadius: sizeImg,
            }} 
            />;
}
