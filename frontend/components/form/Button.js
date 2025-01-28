import { StyleSheet, View, Pressable, Text } from 'react-native';
import FontAwesome from "@expo/vector-icons/FontAwesome";
import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import MaterialCommunityIcons from "@expo/vector-icons/MaterialCommunityIcons";

export default function Button({ label, theme, onPress }) {
    if (theme === "primary") {
        return (
          <View style={[styles.buttonContainer, { borderWidth: 2, borderColor: "#73fc03", borderRadius: 18 }]}>
            <Pressable
              style={[styles.button, { backgroundColor: "#fff" }]}
              onPress={onPress}
            >
              <FontAwesome
                name="picture-o"
                size={18}
                color="#25292e"
                style={styles.buttonIcon}
              />
              <Text style={[styles.buttonLabel, { color: "#25292e" }]}>{label}</Text>
            </Pressable>
          </View>
        );
      }

    if (theme === "save") {
        return (
            <View style={[styles.buttonContainer, { borderWidth: 2, borderColor: "#73fc03", borderRadius: 18 }]}>
                <Pressable
                    style={[styles.button, { backgroundColor: "#fff" }]}
                    onPress={onPress}
                >
                    <MaterialIcons
                        name="save"
                        size={18}
                        color="#25292e"
                        style={styles.buttonIcon}
                    />
                    <Text style={[styles.buttonLabel, { color: "#25292e" }]}>{label}</Text>
                </Pressable>
            </View>
        );
    }

    if (theme === "reset") {
        return (
            <View style={[styles.buttonContainer, { borderWidth: 2, borderColor: "#73fc03", borderRadius: 18 }]}>
                <Pressable
                    style={[styles.button, { backgroundColor: "#fff" }]}
                    onPress={onPress}
                >
                    <MaterialIcons
                        name="restart-alt"
                        size={18}
                        color="#25292e"
                        style={styles.buttonIcon}
                    />
                    <Text style={[styles.buttonLabel, { color: "#25292e" }]}>{label}</Text>
                </Pressable>
            </View>
        );
    }

    if (theme === "send") {
        return (
            <View style={[styles.buttonContainer, { borderWidth: 2, borderColor: "#73fc03", borderRadius: 18 }]}>
                <Pressable
                    style={[styles.button, { backgroundColor: "#fff" }]}
                    onPress={onPress}
                >
                    <MaterialCommunityIcons
                        name="send"
                        size={18}
                        color="#25292e"
                        style={styles.buttonIcon}
                    />
                    <Text style={[styles.buttonLabel, { color: "#25292e" }]}>{label}</Text>
                </Pressable>
            </View>
        );
    }

    return (
        <View style={styles.buttonContainer}>
            <Pressable style={styles.button} onPress={onPress} >
                <Text style={styles.buttonLabel}>{label}</Text>
            </Pressable>
        </View>
    );
}

const styles = StyleSheet.create({
  buttonContainer: {
    width: 130,
    height: 30,
    marginHorizontal: 20,
    padding: 3,
    // marginBottom: 10,
  },
  button: {
    borderRadius: 10,
    width: '100%',
    height: '100%',
    paddingRight: 3,
    flexDirection: 'row',
  },
  buttonIcon: {
    paddingRight: 8,
    paddingLeft: 7,
    paddingTop: 1,
  },
  buttonLabel: {
    color: '#73fc03',
    fontSize: 16,
    top: -2,
    right: 2
  },
});