import { mount } from 'svelte';
import Root from './Root.svelte';
import './ui.css';
import './fonts.css';
import './style.css';

mount(Root, { target: document.getElementById('app')! });
