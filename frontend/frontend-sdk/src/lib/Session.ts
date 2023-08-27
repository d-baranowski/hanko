import { DecodedToken, SessionDetail } from "./events/CustomEvents";
import { SessionState } from "./state/session/SessionState";
import { AuthTokenPersistence } from "./AuthTokenPersistence";
import jwtDecode from "jwt-decode";

/**
 * Options for Session
 *
 * @category SDK
 * @subcategory Internal
 * @property {string} localStorageKey - The prefix / name of the local storage keys.
 */
interface SessionOptions {
  localStorageKey: string;
}

/**
 A class representing a session.

 @category SDK
 @subcategory Session
 @param {SessionOptions} options - The options that can be used
 */
export class Session {
  _sessionState: SessionState;
  jwt: AuthTokenPersistence;

  // eslint-disable-next-line require-jsdoc
  constructor(options: SessionOptions) {
    this._sessionState = new SessionState({ ...options });
    this.jwt = new AuthTokenPersistence({ storageKey: options.localStorageKey });
  }

  /**
   Checks if the user is logged in.

   @returns {boolean} true if the user is logged in, false otherwise.
   */
  public isValid(): boolean {
    const session = this.get();
    return Session.validate(session);
  }

  public isLoggedIn(): boolean {
    return this.isValid()
  }

  /**
   Retrieves the session details.

   @ignore
   @returns {SessionDetail} The session details.
   */
  public get(): SessionDetail | null {
    this._sessionState.read();

    const userID = this._sessionState.getUserID();
    const expirationSeconds = this._sessionState.getExpirationSeconds();
    const storedToken = this.jwt.getStoredToken();
    const decoded = storedToken ? jwtDecode(storedToken) as DecodedToken : null

    const detail = {
      userID,
      expirationSeconds,
      jwt: storedToken,
      decodedJwt: decoded
    };

    return Session.validate(detail) ? detail : null;
  }

  /**
   Checks if the auth flow is completed. The value resets after the next login attempt.

   @returns {boolean} Returns true if the authentication flow is completed, false otherwise
   */
  public isAuthFlowCompleted() {
    this._sessionState.read();
    return this._sessionState.getAuthFlowCompleted();
  }

  /**
   Validates the session.

   @private
   @param {SessionDetail} detail - The session details to validate.
   @returns {boolean} true if the session details are valid, false otherwise.
   */
  private static validate(detail: SessionDetail | null): boolean {
    if (!detail) {
      return false;
    }

    if (detail.expirationSeconds <= 0) {
      return false
    }

    if (!detail.userID?.length) {
      return false
    }

    if (!detail.jwt?.length) {
      return false;
    }
      
    return true
  }
}
