# Namespaces Report

> This is an experiment and Okteto does not officially support it.

This is a simple example on how you can use Okteto's public API to generate a CSV report of the namespaces that Okteto manages.


- Create an [Okteto Admin Token](https://www.okteto.com/docs/admin/dashboard/#admin-access-tokens)

- Export the token to a local variable:

```bash
export OKTETO_TOKEN=<<your-token>>
```

- Export the URL of your Okteto instance to a local variable:

```bash
export OKTETO_URL=<<your-okteto-url>>
```

- Clone the repository:
```
git clone https://github.com/okteto-community/namespaces-report-csv
```

- Build the binary:

```
make build
```

- Run the command:
```
./namespaces
```

Once finished, the program will generate a CSV with the information about your namespaces, and will print out a summary via standard output.
