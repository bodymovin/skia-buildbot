/**
 * @module module/cluster-digests-sk
 * @description <h2><code>cluster-digests-sk</code></h2>
 *
 * This element renders a list of digests in a D3 force layout. That is, it draws them as circles
 * (aka nodes) as if they were attached via springs that are proportional to the difference between
 * the digest images; the nodes repel each other, as if they were charged particles.
 *
 * It is strongly recommended to have the d3 documentation handy, as this element makes heavy use
 * of that (somewhat dense) library.
 *
 * TODO(kjlubick) make this interactive, like the old Polymer element was.
 *
 * @evt layout-complete; fired when the force layout has stabilized (i.e. finished rendering).
 *
 * @evt selection-changed; fired when a digest is clicked on (or the selection is cleared). Detail
 *   contains a list of digests that are selected.
 *
 */
import { define } from 'elements-sk/define';
import { html } from 'lit-html';
import { $$ } from 'common-sk/modules/dom';
import * as d3Force from 'd3-force';
import * as d3Select from 'd3-selection';
import { ElementSk } from '../../../infra-sk/modules/ElementSk';

const template = (ele) => html`
<svg width=${ele._width} height=${ele._height}></svg>
`;

define('cluster-digests-sk', class extends ElementSk {
  constructor() {
    super(template);

    this._nodes = [];
    this._links = [];

    this._linkTightness = 1 / 8;
    this._nodeRepulsion = 256;

    this._width = 400;
    this._height = 400;

    // An array of digests (strings) that correspond to the currently selected digests (if any).
    this._selectedDigests = [];
  }

  connectedCallback() {
    super.connectedCallback();
    this._render();
  }

  changeLinkTightness(isScaleUp) {
    if (isScaleUp) {
      this._linkTightness *= 1.5;
    } else {
      this._linkTightness /= 1.5;
    }
    this._layout();
  }

  changeNodeRepulsion(isScaleUp) {
    if (isScaleUp) {
      this._nodeRepulsion *= 1.5;
    } else {
      this._nodeRepulsion /= 1.5;
    }
    this._layout();
  }

  /**
   * Recomputes the positions of the digest nodes given the value of links. It expects all the SVG
   * elements (e.g. circles, lines) to already be created; this function will simply update the
   * X and Y values accordingly.
   */
  _layout() {
    const clusterSk = $$('svg', this);

    // This force acts as a repulsion force between digest nodes. This acts a lot like charged
    // particles repelling one another. The main purpose here is to keep nodes from overlapping.
    // See https://github.com/d3/d3-force#forceManyBody
    const chargeForce = d3Force.forceManyBody()
      .strength(-this._nodeRepulsion)
      // Given our nodes have a radius of 12, if two nodes are 60 pixels apart, they are definitely
      // not overlapping, so we can stop counting their "charge". This should help performance by
      // reducing computation needs.
      .distanceMax(60);

    // This force acts as a spring force between digest nodes. More similar digests pull more
    // tightly and should be closer together.
    // See https://github.com/d3/d3-force#links
    const linkForce = d3Force.forceLink(this._links)
      .distance((d) => d.value / this._linkTightness);

    // This force keeps the diagram centered in the SVG.
    // See https://github.com/d3/d3-force#centering
    const centerForce = d3Force.forceCenter(this._width / 2, this._height / 2);

    // These forces help keep the nodes in the visible area.
    const xForce = d3Force.forceX(this._width / 2);
    xForce.strength(0.1);
    const yForce = d3Force.forceY(this._height / 2);
    yForce.strength(0.2); // slightly stronger force down since we have more width to draw into

    // This starts a simulation that will render over the next few seconds as the nodes are
    // simulated into place.
    // See https://github.com/d3/d3-force#forceSimulation
    d3Force.forceSimulation(this._nodes)
      .force('charge', chargeForce) // The names are arbitrary (and inspired by D3 documentation).
      .force('link', linkForce)
      .force('center', centerForce)
      .force('fitX', xForce)
      .force('fixY', yForce)
      .alphaDecay(0.03395) // 1 - pow(0.001, 1 / 200); i.e. 200 iterations
      .on('tick', () => {
        // On each tick, the simulation will update the x,y values of the nodes. We can then
        // select and update those nodes.
        d3Select.select(clusterSk)
          .selectAll('.node')
          .attr('cx', (d) => d.x)
          .attr('cy', (d) => d.y);

        d3Select.select(clusterSk)
          .selectAll('.label')
          .attr('x', (d) => d.x + 14) // offset the labels from the center of the nodes.
          .attr('y', (d) => d.y + 20);

        // source and target are supplied and updated by forceLink:
        // https://github.com/d3/d3-force#link_links
        d3Select.select(clusterSk)
          .selectAll('.link')
          .attr('x1', (d) => d.source.x)
          .attr('y1', (d) => d.source.y)
          .attr('x2', (d) => d.target.x)
          .attr('y2', (d) => d.target.y);
      })
      .on('end', () => {
        this.dispatchEvent(new CustomEvent('layout-complete', { bubbles: true }));
      });
  }

  _getNodeCSSClass(d) {
    let base = `node ${d.status}`;
    if (this._selectedDigests.indexOf(d.name) >= 0) {
      base += ' selected';
    }
    return base;
  }

  /**
   * Sets the new data to render in a cluster view.
   *
   * @param nodes Array<Object>: contains Strings for keys digest (called name for historical
   *   reasons), status, and label.
   * @param links Array<Object>: contains Numbers for keys source, target, and value. source and
   *   target refer to the index of the nodes array. value represents how far apart those two
   *   nodes should be.
   */
  setData(nodes, links) {
    this._nodes = nodes;
    this._links = links;

    this._render();
    // For reasons unknown, after render, we don't always see the SVG element rendered in our
    // DOM, so we schedule the drawing for the next animation frame (when we *do* see the SVG
    // in the DOM).
    window.requestAnimationFrame(() => {
      const clusterSk = $$('svg', this);

      // Delete existing SVG elements
      d3Select.select(clusterSk)
        .selectAll('.link,.node,.label')
        .remove();

      // Reset selection.
      this._selectedDigests = [];

      // We don't have any lines or dots spawn in or dynamically get removed from the drawing, so
      // we don't need to supply an id function to the data calls below.

      // Draw the lines first so they are behind the circles.
      d3Select.select(clusterSk)
        .selectAll('line.link')
        .data(this._links)
        .enter()
        .append('line')
        .attr('class', 'link')
        .attr('stroke', '#ccc')
        .attr('stroke-width', '2');

      // Draw the labels behind the circles because the circles are clickable.
      d3Select.select(clusterSk)
        .selectAll('text.label')
        .data(this._nodes)
        .enter()
        .append('text')
        .attr('class', 'label');
      d3Select.select(clusterSk) // update all nodes with the correct label.
        .selectAll('text.label')
        .text((d) => d.label || '');

      d3Select.select(clusterSk)
        .selectAll('circle.node')
        .data(this._nodes)
        .enter()
        .append('circle')
        .attr('class', (d) => this._getNodeCSSClass(d))
        .attr('r', 12)
        .attr('stroke', 'black')
        .attr('data-digest', (d) => d.name)
        .on('click tap', (d) => {
          // Capture this event (prevent it from propagating to the SVG).
          const evt = d3Select.event;
          evt.preventDefault();
          evt.stopPropagation();

          const digest = d.name;
          if (this._selectedDigests.indexOf(digest) >= 0) {
            return; // It's already selected, do nothing.
          }
          if (evt.shiftKey || evt.ctrlKey || evt.metaKey) {
            // Support multiselection if shift, control or meta is held.
            this._selectedDigests.push(digest);
          } else {
            // Clear the existing selection and replace it with this digest.
            this._selectedDigests = [digest];
          }
          this._updateSelection();
        });

      d3Select.select(clusterSk).on('click tap', () => {
        // Capture this event (prevent it from propagating outside the SVG).
        const evt = d3Select.event;
        evt.preventDefault();
        evt.stopPropagation();

        this._selectedDigests = [];
        this._updateSelection();
      });

      this._layout();
    });
  }

  setWidth(w) {
    if (w === this._width) {
      // Don't need to re-render if the width is unchanged.
      return;
    }
    this._width = w;
    this._layout();
  }

  _updateSelection() {
    d3Select.select($$('svg', this))
      .selectAll('circle.node')
      .data(this._nodes)
      .attr('class', (d) => this._getNodeCSSClass(d));

    this.dispatchEvent(new CustomEvent('selection-changed', {
      bubbles: true,
      detail: this._selectedDigests,
    }));
  }
});
