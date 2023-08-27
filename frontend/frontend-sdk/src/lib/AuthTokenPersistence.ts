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
    return (
      window.localStorage.getItem(this.key + "_auth_token-persistence") || ""
    );
  }

  setStoredToken(token: string, secure = true): void {
    window.localStorage.setItem(this.key + "_auth_token-persistence", token);
  }

  removeStoredToken(): void {
    window.localStorage.removeItem(this.key + "_auth_token-persistence");
  }
}
