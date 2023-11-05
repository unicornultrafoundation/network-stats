import { BrowserRouter, Route, Switch } from 'react-router-dom';
import Home from "../pages/Home";
import { Layout } from './Layout';

export function Routing() {

    return (
    // @ts-ignore
    <BrowserRouter>
      <Layout>
        {/*@ts-ignore*/}
        <Switch>
          {/*@ts-ignore*/}
          <Route exact path="/">
            <Home />
          </Route>
        </Switch>
      </Layout>
    </BrowserRouter>
  )
}
