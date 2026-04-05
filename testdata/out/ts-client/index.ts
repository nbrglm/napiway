import * as sdk from "ts-sdk";

/** VALID AUTH */
var VALID = "valid"

/** INVALID AUTH */
var INVALID = "invalid"

var PAGE_SIZE = 20
var AGE = 30
var AGE1 = 28

var users: sdk.User[] = [
  {
    UserId: "1",
    UserName: "Alice",
    Email: "alice@example.com",
    IsActive: true,
    Age: AGE1,
  },
  {
    UserId: "2",
    UserName: "Bob",
    Email: "bob@example.com",
    IsActive: false,
  },
]

const results: Record<string, boolean> = {};

async function runTests() {
  const serverAddr = process.argv[2];

  if (!serverAddr) {
    console.error("Error: Please provide the server address as the first argument.")
  }

  const api = new sdk.TestingAPI(serverAddr);

  try {
    await testUserLogout(api);

    await testListUsers(api);

    await testGetUser(api);

    await testCreateUser(api);

    await testWhoAmI(api);

    // print the results
    console.log(JSON.stringify(results, null, 2));
  } catch (e) {
    console.error("Error: Global test failure: ", e);
    process.exit(1)
  }
}

async function testUserLogout(api: sdk.TestingAPI) {
  try {
    var noApiKeyReq: sdk.LogoutUserReq = {
      APIKeyAuth: "",
    }
    await api.LogoutUser(noApiKeyReq);
  } catch (e) {
    if (e instanceof sdk.TestingAPIError) {
      if (e.reason === sdk.ReasonEncoding) {
        results["LogoutUserWithMissingTokens"] = true;
      }
    } else {
      throw e;
    }
  }
  if (!results["LogoutUserWithMissingTokens"])
    // If the entry isn't in there, mark it false explicitly
    results["LogoutUserWithMissingTokens"] = false;


  try {
    var validApiKeyReq: sdk.LogoutUserReq = {
      APIKeyAuth: VALID
    }
    await api.LogoutUser(validApiKeyReq)
  } catch (e) {
    if (e instanceof sdk.TestingAPIError)
      results["LogoutUserWithValidAPIKey"] = true
    else
      throw e;

  }
  if (!results["LogoutUserWithValidAPIKey"])
    results["LogoutUserWithValidAPIKey"] = false;


  var invalidApiKeyReq: sdk.LogoutUserReq = {
    APIKeyAuth: INVALID,
    SessionTokenAuth: VALID,
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r1 = await api.LogoutUser(invalidApiKeyReq)
  if (r1.StatusCode == 400)
    results["LogoutUserWithInvalidAPIKey"] = true;
  else
    results["LogoutUserWithInvalidAPIKey"] = false;


  var invalidRefreshTokenReq: sdk.LogoutUserReq = {
    APIKeyAuth: VALID,
    RefreshTokenAuth: INVALID,
  }
  const r2 = await api.LogoutUser(invalidRefreshTokenReq);
  if (r2.StatusCode == 400)
    results["LogoutUserWithInvalidRefreshToken"] = true
  else
    results["LogoutUserWithInvalidRefreshToken"] = false


  var invalidSessionTokenReq: sdk.LogoutUserReq = {
    APIKeyAuth: VALID,
    SessionTokenAuth: INVALID,
  }
  const r3 = await api.LogoutUser(invalidSessionTokenReq)
  if (r3.StatusCode == 400)
    results["LogoutUserWithInvalidSessionToken"] = true
  else
    results["LogoutUserWithInvalidSessionToken"] = false

  var validSessionTokenReq: sdk.LogoutUserReq = {
    APIKeyAuth: VALID,
    SessionTokenAuth: VALID,
  }
  const r4 = await api.LogoutUser(validSessionTokenReq)
  if (r4.StatusCode == 200)
    results["LogoutUserWithValidSessionToken"] = true
  else
    results["LogoutUserWithValidSessionToken"] = false

  var validRefreshTokenReq: sdk.LogoutUserReq = {
    APIKeyAuth: VALID,
    RefreshTokenAuth: VALID,
  }
  const r5 = await api.LogoutUser(validRefreshTokenReq)
  if (r5.StatusCode == 200)
    results["LogoutUserWithValidRefreshToken"] = true
  else
    results["LogoutUserWithValidRefreshToken"] = false
}

async function testListUsers(api: sdk.TestingAPI) {
  try {
    var validApiKeyReq: sdk.ListUsersReq = {
      APIKeyAuth: "",
      AdminTokenAuth: "",
    }
    await api.ListUsers(validApiKeyReq);
  } catch (e) {
    if (e instanceof sdk.TestingAPIError) {
      if (e.reason === sdk.ReasonEncoding) {
        results["ListUsersWithMissingAPIKey"] = true;
      }
    } else {
      throw e;
    }
  }
  if (!results["ListUsersWithMissingAPIKey"])
    // If the entry isn't in there, mark it false explicitly
    results["ListUsersWithMissingAPIKey"] = false;

  try {
    var validApiKeyReq: sdk.ListUsersReq = {
      APIKeyAuth: VALID,
      AdminTokenAuth: "",
    }
    await api.ListUsers(validApiKeyReq);
  } catch (e) {
    if (e instanceof sdk.TestingAPIError) {
      if (e.reason === sdk.ReasonEncoding) {
        results["ListUsersWithMissingAdminToken"] = true;
      }
    } else {
      throw e;
    }
  }
  if (!results["ListUsersWithMissingAdminToken"])
    // If the entry isn't in there, mark it false explicitly
    results["ListUsersWithMissingAdminToken"] = false;

  var invalidApiKeyReq: sdk.ListUsersReq = {
    APIKeyAuth: INVALID,
    AdminTokenAuth: VALID,
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r1 = await api.ListUsers(invalidApiKeyReq)
  if (r1.StatusCode == 400)
    results["ListUsersWithInvalidAPIKey"] = true;
  else
    results["ListUsersWithInvalidAPIKey"] = false;

  var invalidAdminTokenReq: sdk.ListUsersReq = {
    APIKeyAuth: VALID,
    AdminTokenAuth: INVALID,
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r2 = await api.ListUsers(invalidAdminTokenReq)
  if (r2.StatusCode == 400)
    results["ListUsersWithInvalidAdminToken"] = true;
  else
    results["ListUsersWithInvalidAdminToken"] = false;

  var validWithoutQueryParamsReq: sdk.ListUsersReq = {
    APIKeyAuth: VALID,
    AdminTokenAuth: VALID,
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r3 = await api.ListUsers(validWithoutQueryParamsReq)
  if (r3.StatusCode == 200)
    results["ListUsersValidOperationWithoutQueryParams"] = true;
  else
    results["ListUsersValidOperationWithoutQueryParams"] = false;

  var validWithQueryParamsReq: sdk.ListUsersReq = {
    APIKeyAuth: VALID,
    AdminTokenAuth: VALID,
    PageSize: PAGE_SIZE,
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r4 = await api.ListUsers(validWithQueryParamsReq)
  if (r4.StatusCode == 200)
    results["ListUsersValidOperationWithQueryParams"] = true;
  else
    results["ListUsersValidOperationWithQueryParams"] = false;
}

async function testGetUser(api: sdk.TestingAPI) {
  try {
    var noApiKeyReq: sdk.GetUserReq = {
      UserId: "",
      APIKeyAuth: VALID,
      SessionTokenAuth: VALID,
    }
    await api.GetUser(noApiKeyReq);
  } catch (e) {
    if (e instanceof sdk.TestingAPIError) {
      if (e.reason === sdk.ReasonEncoding) {
        results["GetUserWithoutQueryParam"] = true;
      }
    } else {
      throw e;
    }
  }
  if (!results["GetUserWithoutQueryParam"])
    // If the entry isn't in there, mark it false explicitly
    results["GetUserWithoutQueryParam"] = false;

  var validReq: sdk.GetUserReq = {
    APIKeyAuth: VALID,
    SessionTokenAuth: VALID,
    UserId: "1",
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r1 = await api.GetUser(validReq)
  if (r1.StatusCode == 200 && r1.Response200.Body.UserId == "1")
    results["GetUserValidOperation"] = true;
  else
    results["GetUserValidOperation"] = false;
}

async function testCreateUser(api: sdk.TestingAPI) {
  var validWithoutOptionalReq: sdk.CreateUserReq = {
    APIKeyAuth: VALID,
    AdminTokenAuth: VALID,
    Body: sdk.createCreateUserRequestBody(
      {
        UserName: "Test User",
        Email: "test@example.com",
        Status: sdk.UserStatusACTIVE,
      },
    )
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r1 = await api.CreateUser(validWithoutOptionalReq)
  if (r1.StatusCode == 201 && r1.Response201.Body.User.Email == "test@example.com" && r1.Response201.Body.Status == sdk.UserStatusACTIVE)
    results["CreateUserValidOperationWithoutOptionalField"] = true;
  else
    results["CreateUserValidOperationWithoutOptionalField"] = false;

  var validWithOptionalReq: sdk.CreateUserReq = {
    APIKeyAuth: VALID,
    AdminTokenAuth: VALID,
    Body: sdk.createCreateUserRequestBody(
      {
        UserName: "Test User",
        Email: "test@example.com",
        Age: AGE,
        Status: sdk.UserStatusACTIVE,
        OptionalStatus: sdk.UserStatusINACTIVE_USER,
      },
    )
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r2 = await api.CreateUser(validWithOptionalReq)
  if (r2.StatusCode == 201 && r2.Response201.Body.User.Age == AGE && r2.Response201.Body.OptionalStatus == sdk.UserStatusINACTIVE_USER)
    results["CreateUserValidOperationWithOptionalField"] = true;
  else
    results["CreateUserValidOperationWithOptionalField"] = false;

  var arbitraryData: Record<string, any> = { "key1": "value1", "key2": 123.0, "key3": true }
  var validWithArbitraryDataReq: sdk.CreateUserReq = {
    APIKeyAuth: VALID,
    AdminTokenAuth: VALID,
    Body: sdk.createCreateUserRequestBody(
      {
        UserName: "Test User",
        Email: "test@example.com",
        Status: sdk.UserStatusACTIVE,
        ArbitraryData: arbitraryData,
      },
    )
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r3 = await api.CreateUser(validWithArbitraryDataReq)
  if (r3.StatusCode == 201 && JSON.stringify(r3.Response201.Body.ArbitraryData) == JSON.stringify(arbitraryData))
    results["CreateUserValidOperationWithArbitraryData"] = true;
  else
    results["CreateUserValidOperationWithArbitraryData"] = false;
}

async function testWhoAmI(api: sdk.TestingAPI) {
  var userId = "test@example.com"
  var validWithoutOptionalReq: sdk.WhoAmIReq = {
    APIKeyAuth: VALID,
    SessionTokenAuth: VALID,
  }
  // No need for a try/catch as an error isn't expected for this case to pass
  const r1 = await api.WhoAmI(validWithoutOptionalReq, userId)
  if (r1.StatusCode == 200 && (await r1.Response200.RawBody.text()) == "test@example.com")
    results["WhoAmIValidRawBody"] = true;
  else
    results["WhoAmIValidRawBody"] = false;
}

// Run the test runner
runTests();
