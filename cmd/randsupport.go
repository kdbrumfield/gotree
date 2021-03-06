package cmd

import (
	goio "io"
	"os"

	"github.com/fredericlemoine/gostats"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

// randsupportCmd represents the randbrlen command
var randsupportCmd = &cobra.Command{
	Use:   "setrand",
	Short: "Assign a random support to edges of input trees",
	Long: `Assign a random support to edges of input trees.

Support follows a uniform distribution in [0,1].

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var f *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees

		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		for tr := range treechan {
			if tr.Err != nil {
				io.LogError(tr.Err)
				return tr.Err
			}
			for _, e := range tr.Tree.Edges() {
				if !e.Right().Tip() {
					e.SetSupport(gostats.Float64Range(0, 1))
				}
			}
			f.WriteString(tr.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	supportCmd.AddCommand(randsupportCmd)

	randsupportCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	randsupportCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
}
