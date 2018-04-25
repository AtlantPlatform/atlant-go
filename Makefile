all:

clean:
	rm -rf var/fs var/state

clean-cluster2: clean-node1 clean-node2
init-cluster2: init-node1 init-node2

clean-cluster3: clean-node1 clean-node2 clean-node3
init-cluster3: init-node1 init-node2 init-node3

init-node1:
	atlant-go -T -F var/fs1 -S var/state1 init

init-node2:
	atlant-go -T -F var/fs2 -S var/state2 init

init-node3:
	atlant-go -T -F var/fs3 -S var/state3 init

clean-node1:
	rm -rf var/fs1 var/state1

clean-node2:
	rm -rf var/fs2 var/state2

clean-node3:
	rm -rf var/fs3 var/state3

run-node1: export IPFS_LOGGING = warning
run-node1:
	atlant-go -E 0x0 -F var/fs1 -S var/state1 -L ":33771" -W ":33781" -l 5 $(COMMAND)

run-node2: export IPFS_LOGGING = warning
run-node2:
	atlant-go -E 0x0 -F var/fs2 -S var/state2 -L ":33772" -W ":33782" -l 5 $(COMMAND)

run-node3: export IPFS_LOGGING = warning
run-node3:
	atlant-go -E 0x0 -F var/fs3 -S var/state3 -L ":33773" -W ":33783" -l 5 $(COMMAND)

install:
	go install -tags testing github.com/AtlantPlatform/atlant-go
