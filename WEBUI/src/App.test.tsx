import React from '../node_modules/@types/react';
import { render } from '../node_modules/@types/testing-library__react';
import App from './App';

test('renders learn react link', () => {
  const { getByText } = render(<App />);
  const linkElement = getByText(/learn react/i);
  expect(linkElement).toBeInTheDocument();
});
