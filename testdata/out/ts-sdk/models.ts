


/**
 * TestingAPIError is a custom error class that encapsulates different types of errors that can occur when using the TestingAPI client. It includes a reason for the error, a message, and optionally the underlying error object for more detailed debugging information.
 */
export class TestingAPIError extends Error {
  reason: TestingAPIErrorReason;
  error?: Error; // Optional field to hold the original error object, if available

  constructor(reason: TestingAPIErrorReason, message: string, error?: Error) {
    super(message);
    this.reason = reason;
    this.error = error;
    if (error) {
      this.stack = error.stack; // Preserve original stack trace if available
    }
  }
}

export type TestingAPIErrorReason =
  | "transport"
  | "encoding"
  | "unexpected";

/** Network/Timeout */
export const ReasonTransport = "transport";

/** Marshal/Unmarshal */
export const ReasonEncoding = "encoding";

/** Non-Spec Response */
export const ReasonUnexpected = "unexpected";





/**
 * Request body for creating a new user.
 */

export interface CreateUserRequestBody {
  
  
  /**
  * The age of the user to be created.
  * Optional
  * 
  */
  Age?: number;

  
  
  /**
  * An object that can contain any arbitrary data related to the user list response. Just for testing freeForm object support in the generator. This is returned by the server in the response body.
  * Optional
  * 
  */
  ArbitraryData?: Record<string, any>;

  
  
  /**
  * The email address of the user to be created.
  * Required
  *  Must be non-empty
  */
  Email: string;

  
  
  /**
  * The name of the user to be created.
  * Required
  *  Must be non-empty
  */
  UserName: string;

  
}

/**
 * createCreateUserRequestBody creates a new instance of CreateUserRequestBody with required fields as parameters
 */
export function createCreateUserRequestBody(props: CreateUserRequestBody): CreateUserRequestBody {
  return props;
}


/**
 * Successful response containing the created user information.
 */

export interface CreateUserResponseBody {
  
  
  /**
  * An object that can contain any arbitrary data related to the user list response. Just for testing freeForm object support in the generator.
  * Optional
  * 
  */
  ArbitraryData?: Record<string, any>;

  
  
  /**
  * The created user information.
  * Required
  * 
  */
  User: User;

  
}

/**
 * createCreateUserResponseBody creates a new instance of CreateUserResponseBody with required fields as parameters
 */
export function createCreateUserResponseBody(props: CreateUserResponseBody): CreateUserResponseBody {
  return props;
}


/**
 * Standard error response schema.
 */

export interface ErrorResponse {
  
  
  /**
  * A detailed debug message for developers. Only passed if in debug mode.
  * Optional
  * 
  */
  DebugMessage?: string;

  
  
  /**
  * An error message which is user-friendly.
  * Required
  *  Must be non-empty
  */
  ErrorMessage: string;

  
}

/**
 * createErrorResponse creates a new instance of ErrorResponse with required fields as parameters
 */
export function createErrorResponse(props: ErrorResponse): ErrorResponse {
  return props;
}


/**
 * Response body for the HealthCheck endpoint.
 */

export interface HealthCheckResponseBody {
  
  
  /**
  * The health status of the API, typically "OK".
  * Required
  *  Must be non-empty
  */
  Status: string;

  
}

/**
 * createHealthCheckResponseBody creates a new instance of HealthCheckResponseBody with required fields as parameters
 */
export function createHealthCheckResponseBody(props: HealthCheckResponseBody): HealthCheckResponseBody {
  return props;
}


/**
 * Successful response containing a list of users.
 */

export interface ListUsersResponseBody {
  
  
  /**
  * The current page number.
  * Required
  * 
  */
  PageNumber: number;

  
  
  /**
  * The number of items per page.
  * Required
  * 
  */
  PageSize: number;

  
  
  /**
  * The total number of users available.
  * Required
  * 
  */
  TotalCount: number;

  
  
  /**
  * No description provided
  * Required
  * 
  */
  Users: User;

  
}

/**
 * createListUsersResponseBody creates a new instance of ListUsersResponseBody with required fields as parameters
 */
export function createListUsersResponseBody(props: ListUsersResponseBody): ListUsersResponseBody {
  return props;
}


/**
 * Response body for the LogoutUser endpoint.
 */

export interface LogoutUserResponseBody {
  
  
  /**
  * A message confirming successful logout.
  * Required
  *  Must be non-empty
  */
  Message: string;

  
}

/**
 * createLogoutUserResponseBody creates a new instance of LogoutUserResponseBody with required fields as parameters
 */
export function createLogoutUserResponseBody(props: LogoutUserResponseBody): LogoutUserResponseBody {
  return props;
}


/**
 * Response Schema for GetUser endpoint.
 */

export interface User {
  
  
  /**
  * The age of the user.
  * Optional
  * 
  */
  Age?: number;

  
  
  /**
  * The email address of the user.
  * Required
  *  Must be non-empty
  */
  Email: string;

  
  
  /**
  * Indicates whether the user is active.
  * Required
  * 
  */
  IsActive: boolean;

  
  
  /**
  * The unique identifier of the user.
  * Required
  *  Must be non-empty
  */
  UserId: string;

  
  
  /**
  * The name of the user.
  * Required
  *  Must be non-empty
  */
  UserName: string;

  
}

/**
 * createUser creates a new instance of User with required fields as parameters
 */
export function createUser(props: User): User {
  return props;
}





const CreateUserReqHTTPMethod = "POST";
const CreateUserReqRoutePath = "/users/new";


/**
 * Create a new user in the system.
 */

export type CreateUserReq = {




  // Authentication parameters (all required)
  
  /**
  * Required Authentication Method
  * Source: header "X-App-Admin-Token"
  *  Description: Authentication method that denotes an admin token passed in the request header. 
  *  Format (NOT ENFORCED): admin_token 
  */
  AdminTokenAuth: string;
  
  /**
  * Required Authentication Method
  * Source: header "X-App-API-Key"
  *  Description: Authentication method that denotes an API key passed in the request header. 
  *  Format (NOT ENFORCED): api_key 
  */
  APIKeyAuth: string;
  



  /**
  * Request body
  */
  Body: CreateUserRequestBody;

};



export type CreateUser201 = {
  

  
  /**
  * Response body
  */
  Body: CreateUserResponseBody;
  
};

export async function ParseCreateUser201(resp: Response): Promise<CreateUser201> {
  var result = {} as CreateUser201;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as CreateUserResponseBody;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for CreateUser201");
    }
  );
  
}



export type CreateUser400 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseCreateUser400(resp: Response): Promise<CreateUser400> {
  var result = {} as CreateUser400;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for CreateUser400");
    }
  );
  
}



// CreateUser413 response has no headers or body
// Payload Too Large - the request body exceeds the maximum allowed size



export type CreateUser500 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseCreateUser500(resp: Response): Promise<CreateUser500> {
  var result = {} as CreateUser500;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for CreateUser500");
    }
  );
  
}



const GetUserReqHTTPMethod = "GET";
const GetUserReqRoutePath = "/users/{userId}";


/**
 * Retrieve user information by user ID.
 */

export type GetUserReq = {

  /**
  * Source: path parameter "{userId}"
  
  * The unique identifier of the user.
  * 
  * Required
  */
  UserId: string;





  // Authentication parameters (all required)
  
  /**
  * Required Authentication Method
  * Source: header "X-App-API-Key"
  *  Description: Authentication method that denotes an API key passed in the request header. 
  *  Format (NOT ENFORCED): api_key 
  */
  APIKeyAuth: string;
  
  /**
  * Required Authentication Method
  * Source: header "X-App-Session-Token"
  *  Description: Authentication method that denotes a session token passed in the request header. 
  *  Format (NOT ENFORCED): session_token 
  */
  SessionTokenAuth: string;
  



};



export type GetUser200 = {
  

  
  /**
  * Response body
  */
  Body: User;
  
};

export async function ParseGetUser200(resp: Response): Promise<GetUser200> {
  var result = {} as GetUser200;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as User;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for GetUser200");
    }
  );
  
}



export type GetUser400 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseGetUser400(resp: Response): Promise<GetUser400> {
  var result = {} as GetUser400;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for GetUser400");
    }
  );
  
}



export type GetUser404 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseGetUser404(resp: Response): Promise<GetUser404> {
  var result = {} as GetUser404;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for GetUser404");
    }
  );
  
}



export type GetUser500 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseGetUser500(resp: Response): Promise<GetUser500> {
  var result = {} as GetUser500;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for GetUser500");
    }
  );
  
}



const ListUsersReqHTTPMethod = "GET";
const ListUsersReqRoutePath = "/users";


/**
 * List users with optional pagination.
 */

export type ListUsersReq = {


  /**
  * Source: query parameter "page"
  
  * The page number for pagination. Default = 0.
  * 
  * Optional
  */
  PageNumber?: number;


  /**
  * Source: query parameter "pageSize"
  
  * The number of items per page for pagination. Default = 10.
  * 
  * Optional
  */
  PageSize?: number;




  // Authentication parameters (all required)
  
  /**
  * Required Authentication Method
  * Source: header "X-App-Admin-Token"
  *  Description: Authentication method that denotes an admin token passed in the request header. 
  *  Format (NOT ENFORCED): admin_token 
  */
  AdminTokenAuth: string;
  
  /**
  * Required Authentication Method
  * Source: header "X-App-API-Key"
  *  Description: Authentication method that denotes an API key passed in the request header. 
  *  Format (NOT ENFORCED): api_key 
  */
  APIKeyAuth: string;
  



};



export type ListUsers200 = {
  
  /**
  * Source: header parameter "X-RateLimit-Remaining"
  
  * The number of remaining requests allowed in the current rate limit window.
  * 
  * Required
  */
  RateLimitRemaining: number;

  

  
  /**
  * Response body
  */
  Body: ListUsersResponseBody;
  
};

export async function ParseListUsers200(resp: Response): Promise<ListUsers200> {
  var result = {} as ListUsers200;
  
  result.RateLimitRemaining = parseintegerParam(resp.headers.get("X-RateLimit-Remaining"), "header: X-RateLimit-Remaining", true)!;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ListUsersResponseBody;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for ListUsers200");
    }
  );
  
}



export type ListUsers400 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseListUsers400(resp: Response): Promise<ListUsers400> {
  var result = {} as ListUsers400;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for ListUsers400");
    }
  );
  
}



export type ListUsers500 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseListUsers500(resp: Response): Promise<ListUsers500> {
  var result = {} as ListUsers500;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for ListUsers500");
    }
  );
  
}



const LogoutUserReqHTTPMethod = "GET";
const LogoutUserReqRoutePath = "/users/logout";


/**
 * Logout the current user.
 */

export type LogoutUserReq = {




  // Authentication parameters (all required)
  
  /**
  * Required Authentication Method
  * Source: header "X-App-API-Key"
  *  Description: Authentication method that denotes an API key passed in the request header. 
  *  Format (NOT ENFORCED): api_key 
  */
  APIKeyAuth: string;
  


  // Authentication parameters (at least one required)
  
  /**
  * Source: header "X-App-Refresh-Token"
  *  Description: Authentication method that denotes a refresh token passed in the request header. 
  *  Format (NOT ENFORCED): refresh_token 
  */
  RefreshTokenAuth?: string;
  
  /**
  * Source: header "X-App-Session-Token"
  *  Description: Authentication method that denotes a session token passed in the request header. 
  *  Format (NOT ENFORCED): session_token 
  */
  SessionTokenAuth?: string;
  


};



export type LogoutUser200 = {
  

  
  /**
  * Response body
  */
  Body: LogoutUserResponseBody;
  
};

export async function ParseLogoutUser200(resp: Response): Promise<LogoutUser200> {
  var result = {} as LogoutUser200;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as LogoutUserResponseBody;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for LogoutUser200");
    }
  );
  
}



export type LogoutUser400 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseLogoutUser400(resp: Response): Promise<LogoutUser400> {
  var result = {} as LogoutUser400;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for LogoutUser400");
    }
  );
  
}



export type LogoutUser500 = {
  

  
  /**
  * Response body
  */
  Body: ErrorResponse;
  
};

export async function ParseLogoutUser500(resp: Response): Promise<LogoutUser500> {
  var result = {} as LogoutUser500;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as ErrorResponse;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for LogoutUser500");
    }
  );
  
}



const WhoAmIReqHTTPMethod = "POST";
const WhoAmIReqRoutePath = "/users/whoami";


/**
 * Get information about the currently authenticated user. Provide the request body as just a string which is set to the user id.
 */

export type WhoAmIReq = {




  // Authentication parameters (all required)
  
  /**
  * Required Authentication Method
  * Source: header "X-App-API-Key"
  *  Description: Authentication method that denotes an API key passed in the request header. 
  *  Format (NOT ENFORCED): api_key 
  */
  APIKeyAuth: string;
  
  /**
  * Required Authentication Method
  * Source: header "X-App-Session-Token"
  *  Description: Authentication method that denotes a session token passed in the request header. 
  *  Format (NOT ENFORCED): session_token 
  */
  SessionTokenAuth: string;
  



};



export type WhoAmI200 = {
  
  /**
  * Source: header parameter "X-RateLimit-Remaining"
  
  * The number of remaining requests allowed in the current rate limit window.
  * 
  * Required
  */
  RateLimitRemaining: number;

  

  
  /**
  * Raw response */
  RawBody: Response;
  
};

export async function ParseWhoAmI200(resp: Response): Promise<WhoAmI200> {
  var result = {} as WhoAmI200;
  
  result.RateLimitRemaining = parseintegerParam(resp.headers.get("X-RateLimit-Remaining"), "header: X-RateLimit-Remaining", true)!;
  
  
  result.RawBody = resp;
  return result;
  
}



// WhoAmI400 response has no headers or body
// Invalid Request



const HealthCheckReqHTTPMethod = "GET";
const HealthCheckReqRoutePath = "/health";


export type HealthCheckReq = {






};



export type HealthCheck200 = {
  

  
  /**
  * Response body
  */
  Body: HealthCheckResponseBody;
  
};

export async function ParseHealthCheck200(resp: Response): Promise<HealthCheck200> {
  var result = {} as HealthCheck200;
  
  
  // Attempt to parse response body as JSON
  return resp.json().then(
    body => {
      result.Body = body as HealthCheckResponseBody;
      return result;
    },
    err => {
      throw new TestingAPIError(ReasonUnexpected, "Error parsing response body for HealthCheck200");
    }
  );
  
}




function parsedoubleParam(value: string | null, name: string, required: boolean): number | undefined {
  if (required && (value === null || value.trim() === "")) {
    throw new Error(`Missing required number parameter: ${name}`);
  }
  if (value === null || value.trim() === "") {
    return undefined; // return default value
  }
  const parsed = parseFloat(value);
  if (isNaN(parsed)) {
    throw new Error(`Failed to parse number parameter: ${value}`);
  }
  return parsed;
}

function parseintegerParam(value: string | null, name: string, required: boolean): number | undefined {
  if (required && (value === null || value.trim() === "")) {
    throw new Error(`Missing required number parameter: ${name}`);
  }
  if (value === null || value.trim() === "") {
    return undefined; // return default value
  }
  const parsed = parseInt(value);
  if (isNaN(parsed)) {
    throw new Error(`Failed to parse number parameter: ${value}`);
  }
  return parsed;
}

function parsebooleanParam(value: string | null, name: string, required: boolean): boolean | undefined {
  if (required && (value === null || value.trim() === "")) {
    throw new Error(`Missing required boolean parameter: ${name}`);
  }
  if (value === null || value.trim() === "") {
    return undefined; // return default value
  }
  if (value.toLowerCase() === "true") {
    return true;
  } else if (value.toLowerCase() === "false") {
    return false;
  } else {
    throw new Error(`Failed to parse boolean parameter: ${value}`);
  }
}

function parsestringParam(value: string | null, name: string, required: boolean): string | undefined {
  if (required && (value === null || value.trim() === "")) {
    throw new Error(`Missing required string parameter: ${name}`);
  }
  if (value === null || value.trim() === "") {
    return undefined; // return default value
  }
  return value;
}
