import { StyleSheet, Image } from 'react-native';

export default function ImageViewerSettings({ placeholderImageSource, selectedImage }) {
    const imageSource = selectedImage  ? selectedImage : placeholderImageSource;

    if (selectedImage) {
      return <Image source={{uri: selectedImage}} style={styles.image}/>
    }  

    return <Image source={imageSource} style={styles.image} />;
}

const styles = StyleSheet.create({
  image: {
    width: 250,
    height: 250,
    borderRadius: 200,
  },
});