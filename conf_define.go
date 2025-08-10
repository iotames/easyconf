package easyconf

func (cf *Conf) StringVar(pval *string, name string, defval, title string, usage ...string) {
	*pval = defval
	cf.addItem(pval, name, defval, title, usage...)
}

func (cf *Conf) BoolVar(pval *bool, name string, defval bool, title string, usage ...string) {
	*pval = defval
	cf.addItem(pval, name, defval, title, usage...)
}

func (cf *Conf) IntVar(pval *int, name string, defval int, title string, usage ...string) {
	*pval = defval
	cf.addItem(pval, name, defval, title, usage...)
}

func (cf *Conf) StringListVar(pval *[]string, name string, defval []string, title string, usage ...string) {
	*pval = defval
	cf.addItem(pval, name, defval, title, usage...)
}

func (cf *Conf) IntListVar(pval *[]int, name string, defval []int, title string, usage ...string) {
	*pval = defval
	cf.addItem(pval, name, defval, title, usage...)
}

func (cf *Conf) Float64Var(pval *float64, name string, defval float64, title string, usage ...string) {
	*pval = defval
	cf.addItem(pval, name, defval, title, usage...)
}
