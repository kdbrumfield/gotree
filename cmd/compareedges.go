package cmd

import (
	"bytes"
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// compareedgesCmd represents the compareedges command
var compareedgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "Compare edges of a reference tree with another tree",
	Long: `Compare edges of a reference tree with another tree

If the compared tree file contains several trees, it will take the first one only
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var refTree *tree.Tree
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		fmt.Fprintf(os.Stderr, "Reference : %s\n", intreefile)
		fmt.Fprintf(os.Stderr, "Compared  : %s\n", intree2file)

		if refTree, err = readTree(intreefile); err != nil {
			io.LogError(err)
			return
		}
		if err = refTree.ReinitIndexes(); err != nil {
			io.LogError(err)
			return
		}
		names := refTree.SortedTips()

		edges1 := refTree.Edges()

		fmt.Printf("tree\tbrid\tlength\tsupport\tterminal\tdepth\ttopodepth\trightname\tfound")
		fmt.Printf("\ttransfer\ttaxatomove\tcomparednodename\tcomparedlength\tcomparedsupport\tcomparedtopodepth\tcomparedid")

		fmt.Printf("\n")
		if treefile, treechan, err = readTrees(intree2file); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()
		for t2 := range treechan {
			if t2.Err != nil {
				io.LogError(t2.Err)
				return t2.Err
			}

			if err = t2.Tree.ReinitIndexes(); err != nil {
				io.LogError(err)
			}

			if err = refTree.CompareTipIndexes(t2.Tree); err != nil {
				return
			}

			edges2 := t2.Tree.Edges()
			for i, e1 := range edges1 {
				dist, closeedge, speciestoadd, speciestoremove := support.MinTransferDist(e1, refTree, t2.Tree, len(names), edges2, false)
				var nodename string = "-"
				found := (dist == 0)
				comparelength := "N/A"
				comparedsupport := "N/A"
				comparedtopodepth := -1
				comparedid := -1

				if closeedge != nil {
					nodename = closeedge.Name(t2.Tree.Rooted())
					comparelength = closeedge.LengthString()
					comparedtopodepth, _ = closeedge.TopoDepth()
					comparedid = closeedge.Id()
				}

				fmt.Printf("%d\t%d\t%s\t%t", t2.Id, i, e1.ToStatsString(false), found)
				var movedtaxabuf bytes.Buffer
				for k, sp := range speciestoadd {
					if k > 0 {
						movedtaxabuf.WriteRune(',')
					}
					movedtaxabuf.WriteRune('+')
					movedtaxabuf.WriteString(sp.Name())
				}
				for k, sp := range speciestoremove {
					if k > 0 || (k == 0 && len(speciestoadd) > 0) {
						movedtaxabuf.WriteRune(',')
					}
					movedtaxabuf.WriteRune('-')
					movedtaxabuf.WriteString(sp.Name())
				}
				fmt.Printf("\t%d\t%s\t%s\t%s\t%s\t%d\t%d\n", dist, movedtaxabuf.String(), nodename, comparelength, comparedsupport, comparedtopodepth, comparedid)
			}
		}
		return
	},
}

func init() {
	compareCmd.AddCommand(compareedgesCmd)
	compareedgesCmd.PersistentFlags().BoolVarP(&transferdist, "transfer-dist", "m", false, "If transfer dist must be computed for each edge")
	compareedgesCmd.PersistentFlags().BoolVar(&movedtaxa, "moved-taxa", false, "only if --transfer-dist is given: Then display, for each branch, taxa that must be moved")
}
