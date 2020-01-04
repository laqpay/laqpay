export const AppConfig = {
  otcEnabled: false,
  maxHardwareWalletAddresses: 1,
  useHwWalletDaemon: true,
  urlForHwWalletVersionChecking: 'https://version.laqpay.net/laqwallet/version.txt',
  hwWalletDownloadUrlAndPrefix: 'https://downloads.laqpay.net/laqwallet/laqwallet-firmware-v',

  urlForVersionChecking: 'https://version.laqpay.net/laqpay/version.txt',
  walletDownloadUrl: 'https://www.laqpay.net/downloads/',

  /**
   * This wallet uses the Laqpay URI Specification (based on BIP-21) when creating QR codes and
   * requesting coins. This variable defines the prefix that will be used for creating QR codes
   * and URLs. IT MUST BE UNIQUE FOR EACH COIN.
   */
  uriSpecificatioPrefix: 'laqpay',

  languages: [{
      code: 'en',
      name: 'English',
      iconName: 'en.png',
    },
    {
      code: 'zh',
      name: '中文',
      iconName: 'zh.png',
    },
    {
      code: 'es',
      name: 'Español',
      iconName: 'es.png',
    },
  ],
  defaultLanguage: 'en',
};
