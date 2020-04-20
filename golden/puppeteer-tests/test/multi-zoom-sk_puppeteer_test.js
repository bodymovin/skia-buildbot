const expect = require('chai').expect;
const setUpPuppeteerAndDemoPageServer = require('./util').setUpPuppeteerAndDemoPageServer;
const takeScreenshot = require('./util').takeScreenshot;

describe('multi-zoom-sk', () => {
  const testBed = setUpPuppeteerAndDemoPageServer(); // Contains page and baseUrl.

  beforeEach(async () => {
    await testBed.page.goto(`${testBed.baseUrl}/dist/multi-zoom-sk.html`, { waitUntil: 'networkidle0' });
  });

  it('should render the demo page', async () => {
    // Smoke test. There are 5 test cases we send to Gold and 1 extra in a dialog that we do not.
    expect(await testBed.page.$$('multi-zoom-sk')).to.have.length(6);
  });

  describe('screenshots', () => {
    it('looks good when left and right are a normal size', async () => {
      const multiZoomSk = await testBed.page.$('#normal');
      await takeScreenshot(multiZoomSk, 'multi-zoom-sk');
    });

    it('shows two images of different size', async () => {
      const multiZoomSk = await testBed.page.$('#mismatch');
      await takeScreenshot(multiZoomSk, 'multi-zoom-sk_mismatch');
    });

    it('works for two small base64 encoded images', async () => {
      const multiZoomSk = await testBed.page.$('#base64');
      await takeScreenshot(multiZoomSk, 'multi-zoom-sk_base64-small');
    });

    it('is zoomed in with the grid on', async () => {
      const multiZoomSk = await testBed.page.$('#zoomed_grid');
      await takeScreenshot(multiZoomSk, 'multi-zoom-sk_zoomed-grid');
    });

    it('shows nth largest pixel', async () => {
      const multiZoomSk = await testBed.page.$('#base64_nthpixel');
      await takeScreenshot(multiZoomSk, 'multi-zoom-sk_nth-pixel');
    });
  });
});