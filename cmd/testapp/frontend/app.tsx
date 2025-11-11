import '@molpe/spark/debug';
import { render } from 'preact';
import Client from '@molpe/spark/client';
import type { App } from '@/gotypes';
import { useSignal } from '@preact/signals';
import { ForEach, bindInput } from '@molpe/spark/util';
import '@fortawesome/fontawesome-free/css/all.css';
import 'bootstrap/dist/css/bootstrap.css';
import { WebSocketAsync } from '@molpe/spark/websocket';

export const client = Client<App>('/rpc/')

const ws = new WebSocketAsync('/', { origin: 'wss://echo.websocket.org' });

await ws.connect();

ws.once('text', msg => {
	console.log('from eventlistener', msg.detail);
});

!async function(){
	console.log('begin');

	let ctr = 0;

	for await (const { data, type } of ws) {
		console.log(type, data);
		ctr++;
		if (ctr == 4) {
			break;
		}
	}

	console.log('end');

}();

ws.send('test1');
ws.send('test2');
ws.send('test3');

function App() {

	const text = useSignal('')

	const { data, isLoading } = client.useGetTime(undefined, { refreshInterval: 500, dedupingInterval: 10 });

	if (isLoading) return <span>loading</span>

	return (
		<div class="d-flex h-100 justify-content-center align-items-center flex-column">
			Hello World!
			<div class="mt-2 border rounded p-2">
				RPC Example
				<div class="input-group mb-3">
					<input {...bindInput(text)} type="text" class="form-control" placeholder="Text to Echo"/>
					<button class="btn btn-outline-secondary" onClick={() => client.ExamplePlain({message: text.value}).then(echo => text.value = echo)}>Send!</button>
				</div>
				<span>Server Time {new Date(data?.Time! * 1000).toLocaleString('de-DE')}</span>
			</div>
		</div>
	)
}

render(<App />, document.getElementById('app')!)