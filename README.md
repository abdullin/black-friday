# Python and grpc.


To fix reference problems follow this [link](https://youtrack.jetbrains.com/issue/PY-27111/Please-add-protobuf-autocompletion-support)

You can improve your development experience by installing two Python packages:

https://github.com/dropbox/mypy-protobuf
https://github.com/python/typeshed
$ pip install mypy-protobuf types-protobuf

Once you have these packages installed in your development environment (or whatever environment you use to compile the .proto file if you make certain scripts available on your PATH - please see the docs of mypy-protobuf for details), you can run:

$ protoc mydata.proto --python_out=./ --mypy_out=./

You will get two files created - one .py file with Python classes and one .pyi file with the typing metadata.
Now PyCharm will have Intellisense for your Python code where you access protobuf objects. As a bonus, you can also run mypy on your source code and it will be able to type check it since it has necessary typing metadata to be able to do that.