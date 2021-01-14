/* eslint-disable */
/**
 * @module skottie-text-editor
 * @description <h2><code>skottie-gif-exporter</code></h2>
 *
 * <p>
 *   A skottie gif exporter
 * </p>
 *
 *
 * @evt publish - This event is generated when the user presses Apply.
 *         The updated json is available in the event detail.
 *
 *
 */
import { define } from 'elements-sk/define';
import 'elements-sk/select-sk';
import { html, render } from 'lit-html';
import { ifDefined } from 'lit-html/directives/if-defined.js';
import GIF from './gif'
import GIF_WORKER from './gif.worker'

const QUALITY_SCRUBBER_RANGE = 50;

const ditherOptions = [
  'FloydSteinberg',
  'FalseFloydSteinberg',
  'Stucki',
  'Atkinson'
]

const renderHeader = (ele) => {
  if (ele._state.state === 'idle') {
    return html`
      <button class="editor-header-save-button" @click=${ele._save}>Save</button>
    `
  } else {
    return html`
      <button class="editor-header-save-button" @click=${ele._cancel}>Cancel</button>
    `
  }
}

const renderIdle = ele => {
  return html`
    <div class=form>
      <div class=form-elem>
        <div>Sample</div>
        <input id=sampleScrub type=range min=1 max=${QUALITY_SCRUBBER_RANGE} step=1
            @input=${ele._onSampleScrub} @change=${ele._onSampleScrubEnd}>
      </div>
      <div class=form-elem>
        <checkbox-sk label="Dither"
           ?checked=${ele._state.dither}
           @click=${ele._toggleDither}>
        </checkbox-sk>
        ${
          ele._state.dither
            ? html`
            <select-sk
              role="listbox"
              @selection-changed=${ele._ditherOptionChange}
            >
              ${ditherOptions.map((item, index) => html`
                <div
                  role="option"
                  selected=${ifDefined(ele._state.ditherValue === index ? 'true' : undefined)}
                >
                  ${item}
                </div>
              `)}
            </select-sk>`
            : null
        }
      </div>
      <div class=form-elem>
        <label class=number>
          <input
            type=number
            id=repeats
            .value=${ele._state.repeat}
            min=-1 
            @input=${ele._onRepeatChange}
            @change=${ele._onRepeatChange}
          /> Repeats
        </label>
      </div>
    </div>
    ${ele._state.blob 
      ? html`<a href=${URL.createObjectURL( ele._state.blob )} download>Download</a>`
      : html`<span></span>`
    }
  `
}

const renderImage = ele => {
  return html`
    <div>
      Creating snapshots: ${ele._state.imageProgress}%
    </div>
  `
}

const renderGif = ele => {
  return html`
    <div>
      Creating GIF: ${ele._state.gifProgress}%
    </div>
  `
}

const renderMain = ele => {
  if (ele._state.state === 'idle') {
    return renderIdle(ele);
  } else if (ele._state.state === 'image') {
    return renderImage(ele);
  } else if (ele._state.state === 'gif') {
    return renderGif(ele);
  }
}

const template = (ele) => {
  return html`
  <div>
    <header class="editor-header">
      <div class="editor-header-title">Gif Exporter</div>
      <div class="editor-header-separator"></div>
      ${renderHeader(ele)}
    </header>
    <section>
      ${renderMain(ele)}
    <section>
  </div>
  `;
};

class SkottieGifExporterSk extends HTMLElement {
  constructor() {
    super();
    this._state = {
      workers: 4,
      quality: 10,
      repeat: -1,
      background: '#ffffff',
      dither: false,
      ditherValue: 0,
      state: 'idle',
      gifProgress: 0,
      imageProgress: 0,
      blob: null,
    };
  }

  delay(time) {
    return new Promise((resolve) => setTimeout(resolve, time))
  }

  _onSampleScrub(ev) {
    this._state.quality = ev.target.value;
    this.render();
  }

  _onSampleScrubEnd(ev) {
    this._state.quality = ev.target.value;
    this.render();
  }

  _onRepeatChange(ev) {
    this._state.repeat = ev.target.value;
  }

  _toggleDither(e) {
    e.preventDefault();
    this._state.dither = !this._state.dither;
    this.render();
  }

  _ditherOptionChange(e) {
    e.preventDefault();
    const index = e.detail.selection;
    this._state.ditherValue = e.detail.selection;
    this.render();
  }

  async _processFrames() {
    this.setState
    const fps = this.player.fps();
    const duration = 500;
    // const duration = this.player.duration();
    const canvasElement = this.player.getCanvasElement();
    let currentTime = 0;
    const increment =  1000 / fps;
    this._state.state = 'image';
    this.render();
    while (currentTime < duration) {
      if (this._state.state !== 'image') {
        return;
      }
      await this.delay(1);
      this.player.seekFrame(currentTime / duration);
      this.gif.addFrame(canvasElement, {delay: increment, copy: true});
      this._state.imageProgress = Math.round(currentTime / duration * 100);
      currentTime += increment;
      this.render();
    }
    this._state.state = 'gif';

    this.gif.render();
  }

  _cancel() {
    if (this._state.state === 'gif') {
      this.gif.abort();
    }
    this._state.state = 'idle';
    this.render();
  }

  async _save() {

    this.gif = new GIF({
      workers: this._state.workers,
      quality: this._state.quality,
      repeat: this._state.repeat,
      dither: this._state.dither,
    });
    this.gif.on('finished', blob => {
      // window.open(URL.createObjectURL(blob));
      this._state.state = 'idle';
      this._state.blob = blob;
      this.render();
    });
    this.gif.on('progress', (value) => {
      this._state.gifProgress = Math.round(value * 100);
      this.render();
    });
    this._state.gifProgress = 0;
    this._state.imageProgress = 0;
    this.player.pause();
    this._processFrames();
    
  }

  connectedCallback() {
    // console.log(this.animation);
    this._player = this.player;
    this.render();
    this.addEventListener('input', this._inputEvent);
  }

  disconnectedCallback() {
    this.removeEventListener('input', this._inputEvent);
  }



  render() {
    render(template(this), this, { eventContext: this });
  }
}

define('skottie-gif-exporter', SkottieGifExporterSk);
