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

export const restHandlers = [
  // todo: update when the backend has a real host
  // https://mswjs.io/docs/basics/request-matching#path-with-wildcard
  rest.get('*/rpc/bin_devices_per_hour', (req, res, ctx) => {
    let requestedID = req.url.searchParams.get("_fscs_id");
    let requestedDay = req.url.searchParams.get("_day");
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
]


const server = setupServer(...restHandlers)

// Start server before all tests
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))

//  Close server after all tests
afterAll(() => server.close())

// Reset handlers after each test `important for test isolation`
afterEach(() => server.resetHandlers())