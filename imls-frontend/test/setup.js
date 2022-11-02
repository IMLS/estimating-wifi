import { afterAll, afterEach, beforeAll } from 'vitest'
import { setupServer } from 'msw/node'
import { graphql, rest } from 'msw'
import fetch from 'node-fetch';
// import canvas from 'canvas';
// before enabling a canvas mock, check status on https://github.com/vitest-dev/vitest/issues/740

global.fetch = fetch;
// global.canvas = canvas;

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

const knownGoodDevicesPerHourMock1 = [0,1,2,3,4];
const knownGoodDevicesPerHourMock2 = [9,9,9,9,9];
const knownEmptyDevicesPerHourMock = [0,0,0,0,0];
const errorMock = {
  details: "unexpected \"A\" expecting delimiter (.)",
  message: "\"failed to parse filter (notARealID)\" (line 1, column 4)"
}

const knownGoodLibraryMock = {
  "stabr":"AK",
  "fscskey":"AK0001",
  "fscs_seq":2,
  "libname":"ANCHOR POINT PUBLIC LIBRARY",
  "address":"34020 NORTH FORK ROAD",
  "city":"ANCHOR POINT",
  "zip":"99556",
}

const statewideLibrariesMock = [
  {...knownGoodLibraryMock}
]

export const restHandlers = [
  // todo: update when the backend has a real host
  // https://mswjs.io/docs/basics/request-matching#path-with-wildcard
  rest.get('*/rpc/bin_devices_per_hour', (req, res, ctx) => {
    let requestedID = req.url.searchParams.get("_fscs_id");
    let requestedDay = req.url.searchParams.get("_start");
    switch (requestedID) {
      case 'KnownGoodId' :
        if (requestedDay == "9999-99-99") {
          return res(ctx.status(200), ctx.json(knownGoodDevicesPerHourMock2))
        }
        return res(ctx.status(200), ctx.json(knownGoodDevicesPerHourMock1))
      case 'KnownEmptyId':
        return res(ctx.status(200), ctx.json(knownEmptyDevicesPerHourMock))
      case 'notARealID':
      default:
        return res(ctx.status(400), ctx.json(errorMock))
    }
  }),
  rest.get('*/rpc/bin_devices_over_time', (req, res, ctx) => {
    let requestedID = req.url.searchParams.get("_fscs_id");
    let requestedDay = req.url.searchParams.get("_start"); 
    // todo test these other query params
    let requestedDirection = req.url.searchParams.get("_direction");
    let requestedDuration = req.url.searchParams.get("_days");
    switch (requestedID) {
      case 'KnownGoodId' :
        if (requestedDay == "9999-99-99") {
          return res(ctx.status(200), ctx.json([knownGoodDevicesPerHourMock2, knownGoodDevicesPerHourMock1]))
        }
        return res(ctx.status(200), ctx.json([knownGoodDevicesPerHourMock1, knownGoodDevicesPerHourMock2]))
      case 'KnownEmptyId':
        return res(ctx.status(200), ctx.json([knownEmptyDevicesPerHourMock, knownEmptyDevicesPerHourMock]))
      case 'notARealID':
      default:
        return res(ctx.status(400), ctx.json(errorMock))
    }
  }),
  rest.get('*/rpc/lib_search_fscs', (req, res, ctx) => {
    let requestedID = req.url.searchParams.get("_fscs_id");
    switch (requestedID) {
      case 'KnownGoodId' :
        return res(ctx.status(200), ctx.json(knownGoodLibraryMock))
      default:
        return res(ctx.status(400), ctx.json(errorMock))
    }
  }),
  rest.get('*/rpc/lib_search_state', (req, res, ctx) => {
    let stateAbbr = req.url.searchParams.get("_state_code");
    switch (stateAbbr) {
      case 'AK' :
        return res(ctx.status(200), ctx.json(statewideLibrariesMock))
      case 'AL' :
        return res(ctx.status(200), ctx.json(statewideLibrariesMock))
      case 'ZZ' :
        return res(ctx.status(400), ctx.json(errorMock))
      default:
        return res(ctx.status(400), ctx.json(errorMock))
    }
  }),
  rest.get('*/rpc/lib_search_name', (req, res, ctx) => {
    let textString = req.url.searchParams.get("_name");
    switch (textString) {
      case 'anchor point' :
        return res(ctx.status(200), ctx.json(statewideLibrariesMock))
      default:
        return res(ctx.status(400), ctx.json(errorMock))
    }
  }),
]


const server = setupServer(...restHandlers)

// Start server before all tests
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))

//  Close server after all tests
afterAll(() => server.close())

// Reset handlers after each test `important for test isolation`
afterEach(() => server.resetHandlers())