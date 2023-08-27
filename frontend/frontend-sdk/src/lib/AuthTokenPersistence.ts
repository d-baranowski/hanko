import KeyValueStorageApi from 'src/util/KeyValueStorageApi';

interface AuthTokenPersistanceOptions {
  storageKey: string;
}

export class AuthTokenPersistence {
  key: string;

  // eslint-disable-next-line require-jsdoc
  constructor(options: AuthTokenPersistanceOptions) {
    this.key = options.storageKey;
  }
  getStoredToken(): string {
    return KeyValueStorageApi.getItem(this.key + "_auth_token-persistence") || ""
  }

  setStoredToken(token: string, secure = true): void {
    KeyValueStorageApi.setItem(this.key + "_auth_token-persistence", token)
  }

  removeStoredToken(): void {
    KeyValueStorageApi.removeItem(this.key + "_auth_token-persistence")
  }
}
