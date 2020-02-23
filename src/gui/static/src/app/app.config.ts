export const AppConfig = {
  otcEnabled: true,
  maxHardwareWalletAddresses: 1,
  useHwWalletDaemon: true,
  urlForHwWalletVersionChecking: 'https://api.laqpay.com/wallet/version',
  hwWalletDownloadUrlAndPrefix: 'https://api.laqpay.com/wallet/laq-wallet-firmware-v',

  urlForVersionChecking: 'https://api.laqpay.com/wallet/version',
  walletDownloadUrl: 'https://dl.laqpay.com',

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
