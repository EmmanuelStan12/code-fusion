const ts = require('typescript')

const transpiledCode = ts.transpileModule("interface ICount {\n    value: number;\n}\n\nconst count: ICount = {\n    value: 1,\n};\n\ncount.value++\nconsole.log(count.value, 'Value in console')\ncount.value", {
    compilerOptions: { module: ts.ModuleKind.CommonJS },
}).outputText;
console.log(transpiledCode)