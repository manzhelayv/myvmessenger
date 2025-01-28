// Learn more https://docs.expo.io/guides/customizing-metro
const { getDefaultConfig } = require('expo/metro-config');

/** @type {import('expo/metro-config').MetroConfig} */
const defaultConfig = getDefaultConfig(__dirname);

//defaultConfig.resolver.assetExts.push("cjs");
//defaultConfig.resolver.sourceExts.push(['jsx', 'js', 'ts']);
//defaultConfig.resolver.sourceExts.push('jsx', 'js', 'ts');

module.exports = defaultConfig;

