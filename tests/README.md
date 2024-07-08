Tests in package 'test' should only return an error.
        If the test was successful, return nil.
        If the test failed return ErrTestFailed.
        If there was an error, return err.

func NameOfTest(...prarams) error