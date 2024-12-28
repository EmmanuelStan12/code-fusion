const App = require('./src/app')

const PORT = process.env.PORT || '4000'

const app = new App(PORT)

app.listen()