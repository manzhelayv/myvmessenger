import BackgroundService from 'react-native-background-actions';

const sleep = (time) => new Promise((resolve) => setTimeout(() => resolve(), time));

const veryIntensiveTask = async (taskDataArguments) => {
    const { delay, value } = taskDataArguments;
    await new Promise( async (resolve) => {
        for (let i = 0; BackgroundService.isRunning(); i++) {
            await sleep(10000);
        }
    });
};

async function startBackgroundService(value){
    const options = {
        taskName: 'MYVmessanger',
        taskTitle: 'MYVmessanger',
        taskDesc: 'MYVmessanger в фоновом режиме',
        taskIcon: {
            name: 'ic_launcher',
            type: 'mipmap',
        },
        color: '#ff00ff',
        linkingURI: 'yourSchemeHere://chat/jane', // See Deep Linking for more info
        parameters: {
            delay: 2000,
            value: value,
        },
    };

    await BackgroundService.start(veryIntensiveTask, options);

}

export default function BackgroundNotification(value){
    startBackgroundService(value);
}
