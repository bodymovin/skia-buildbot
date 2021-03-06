/// <reference path="../infra-sk/html-webpack-inject-attributes-plugin.d.ts" />

import { resolve } from 'path';
import webpack from 'webpack';
import CopyPlugin from 'copy-webpack-plugin';
import HtmlWebpackInjectAttributesPlugin from 'html-webpack-inject-attributes-plugin';
import commonBuilder from '../infra-sk/pulito/webpack.common';

const configFactory: webpack.ConfigurationFactory = (_, args) => {
  // Don't minify the HTML since it contains Go template tags.
  const config = commonBuilder(__dirname, args.mode, /* neverMinifyHtml= */ true);

  config.output!.publicPath = '/static/';

  config.plugins!.push(
    new CopyPlugin({
      patterns: [
        {
          from: resolve(__dirname, 'images/icon.png'),
          to: 'icon.png',
        },
        {
          from: resolve(__dirname, 'images/icon-active.png'),
          to: 'icon-active.png',
        },
      ],
    }),
  );

  config.plugins!.push(
    new HtmlWebpackInjectAttributesPlugin({
      nonce: '{%.nonce%}',
    }),
  );

  config.resolve = config.resolve || {};

  // https://github.com/webpack/node-libs-browser/issues/26#issuecomment-267954095
  config.resolve.modules = [resolve(__dirname, 'node_modules')];

  return config;
};

export = configFactory;
