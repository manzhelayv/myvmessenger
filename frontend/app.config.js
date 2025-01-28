export default {
  expo: {
    name: "MYVMessanger",
    slug: "myvmessanger",
    version: "1.0.0",
    orientation: "portrait",
    icon: "./assets/icon.png",
    userInterfaceStyle: "light",
    scheme: "myapp",
    splash: {
      image: "./assets/splash.png",
      resizeMode: "contain",
      backgroundColor: "#25292e"
    },
    ios: {
      supportsTablet: true,
      bundleIdentifier: "com.yvmanzhela.MYVMessanger"
    },
    android: {
      adaptiveIcon: {
        foregroundImage: "./assets/adaptive-icon.png",
        backgroundColor: "#ffffff"
      },
      package: "com.yvmanzhela.MYVMessanger",
      "permissions": [
        "ACCESS_BACKGROUND_LOCATION",
        "ACCESS_COARSE_LOCATION",
        "ACCESS_FINE_LOCATION",
        "FOREGROUND_SERVICE"
      ],
      useNextNotificationsApi: true,
      googleServicesFile: "./google-services.json"
    },
    web: {
      favicon: "./assets/favicon.png"
    },
    extra: {
      eas: {
        projectId: "0f541b61-a9ad-4cfc-a5e0-639d27ec2119"
      }
    },
    plugins:[
      [
        'expo-build-properties',
        {
          android: {
            usesCleartextTraffic: true,
          },
          ios: {
            flipper: true,
          },
        },
      ],
      [
        "expo-secure-store",
        {
          "faceIDPermission": "Allow $(PRODUCT_NAME) to access your Face ID biometric data."
        },
      ],
      [
        "react-native-background-fetch",
      ],
      [
        "./plugins/backgroundTimer/app.plugin.js",
        "backgroundTimer",
      ],
      [
        "expo-contacts",
        {
          "contactsPermission": "Allow $(PRODUCT_NAME) to access your contacts."
        }
      ],
      [
        "expo-font",
        {
          "fonts": ["./assets/fonts/Wittgenstein-Black.otf"]
        }
      ],
      [
        "expo-router"
      ]
      // [
      //   "expo-notifications",
      //   {
      //     "icon": "./assets/favicon.png",
      //     "color": "#191",
      //     "defaultChannel": "default",
      //     // "sounds": [
      //     //   "./local/assets/notification-sound.wav",
      //     //   "./local/assets/notification-sound-other.wav"
      //     // ]
      //   }
      // ]
    ],
  },
}
